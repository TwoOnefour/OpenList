package teldrive

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/OpenListTeam/OpenList/v4/drivers/base"
	"github.com/OpenListTeam/OpenList/v4/internal/driver"
	"github.com/OpenListTeam/OpenList/v4/internal/model"
	"github.com/OpenListTeam/OpenList/v4/internal/stream"
	"github.com/OpenListTeam/OpenList/v4/pkg/utils"
	"github.com/avast/retry-go"
	"github.com/go-resty/resty/v2"
	"golang.org/x/net/context"
)

// create empty file
func (d *Teldrive) touch(name, path string) error {
	uploadBody := base.Json{
		"name": name,
		"type": "file",
		"path": path,
	}
	if err := d.request(http.MethodPost, "/api/files", func(req *resty.Request) {
		req.SetBody(uploadBody)
	}, nil); err != nil {
		return err
	}

	return nil
}

func (d *Teldrive) createFileOnUploadSuccess(name, id, path string, uploadedFileParts []FilePart, totalSize int64) error {
	remoteFileParts, err := d.getFilePart(id)
	if err != nil {
		return err
	}
	// check if the uploaded file parts match the remote file parts
	if len(remoteFileParts) != len(uploadedFileParts) {
		return fmt.Errorf("[Teldrive] file parts count mismatch: expected %d, got %d", len(uploadedFileParts), len(remoteFileParts))
	}
	formatParts := make([]base.Json, 0)
	for _, p := range remoteFileParts {
		formatParts = append(formatParts, base.Json{
			"id":   p.PartId,
			"salt": p.Salt,
		})
	}
	uploadBody := base.Json{
		"name":  name,
		"type":  "file",
		"path":  path,
		"parts": formatParts,
		"size":  totalSize,
	}
	// create file here
	if err := d.request(http.MethodPost, "/api/files", func(req *resty.Request) {
		req.SetBody(uploadBody)
	}, nil); err != nil {
		return err
	}

	return nil
}

func (d *Teldrive) checkFilePartExist(fileId string, partId int) (FilePart, error) {
	var uploadedParts []FilePart
	var filePart FilePart

	if err := d.request(http.MethodGet, "/api/uploads/{id}", func(req *resty.Request) {
		req.SetPathParam("id", fileId)
	}, &uploadedParts); err != nil {
		return filePart, err
	}

	for _, part := range uploadedParts {
		if part.PartId == partId {
			return part, nil
		}
	}

	return filePart, nil
}

func (d *Teldrive) getFilePart(fileId string) ([]FilePart, error) {
	var uploadedParts []FilePart
	if err := d.request(http.MethodGet, "/api/uploads/{id}", func(req *resty.Request) {
		req.SetPathParam("id", fileId)
	}, &uploadedParts); err != nil {
		return nil, err
	}

	return uploadedParts, nil
}

func (d *Teldrive) singleUploadRequest(ctx context.Context, fileId string, callback base.ReqCallback, resp interface{}) error {
	url := d.Address + "/api/uploads/" + fileId
	client := resty.New().SetTimeout(0)

	req := client.R().
		SetContext(ctx)
	req.SetHeader("Cookie", d.Cookie)
	req.SetHeader("Content-Type", "application/octet-stream")
	req.SetContentLength(true)
	req.AddRetryCondition(func(r *resty.Response, err error) bool {
		return false
	})
	if callback != nil {
		callback(req)
	}
	if resp != nil {
		req.SetResult(resp)
	}
	var e ErrResp
	req.SetError(&e)
	_req, err := req.Execute(http.MethodPost, url)
	if err != nil {
		return err
	}

	if _req.IsError() {
		return &e
	}
	return nil
}

func (d *Teldrive) doSingleUpload(ctx context.Context, dstDir model.Obj, file model.FileStreamer, up model.UpdateProgress,
	maxRetried, totalParts int, chunkSize int64, fileId string) error {

	totalSize := file.GetSize()
	var fileParts []FilePart
	var uploaded int64 = 0
	chunkSize = min(totalSize, chunkSize)
	ss, err := stream.NewStreamSectionReader(file, int(chunkSize), &up)
	if err != nil {
		return err
	}
	chunkCnt := 0
	for uploaded < totalSize {
		if utils.IsCanceled(ctx) {
			return ctx.Err()
		}
		curChunkSize := min(totalSize-uploaded, chunkSize)
		rd, err := ss.GetSectionReader(uploaded, curChunkSize)
		if err != nil {
			return err
		}
		filePart := &FilePart{}
		chunkCnt += 1
		if err := retry.Do(func() error {
			if _, err := rd.Seek(0, io.SeekStart); err != nil {
				return err
			}
			if err := d.singleUploadRequest(ctx, fileId, func(req *resty.Request) {
				uploadParams := map[string]string{
					"partName": func() string {
						digits := len(strconv.Itoa(chunkCnt))
						return file.GetName() + fmt.Sprintf(".%0*d", digits, chunkCnt)
					}(),
					"partNo":   strconv.Itoa(chunkCnt),
					"fileName": file.GetName(),
				}
				req.SetQueryParams(uploadParams)
				req.SetBody(driver.NewLimitedUploadStream(ctx, rd))
				req.SetHeader("Content-Length", strconv.FormatInt(curChunkSize, 10))
			}, filePart); err != nil {
				return err
			}
			return nil
		},
			retry.Context(ctx),
			retry.Attempts(uint(maxRetried)),
			retry.DelayType(retry.BackOffDelay),
			retry.Delay(time.Second)); err != nil {
			return err
		}
		if filePart.Name != "" {
			fileParts = append(fileParts, *filePart)
			uploaded += curChunkSize
			up(float64(uploaded) / float64(totalSize) * 100)
			ss.FreeSectionReader(rd)
		} else {
			// For common situation this code won't reach
			return fmt.Errorf("[Teldrive] upload chunk %d failed: filePart Somehow missing", chunkCnt)
		}
	}

	return d.createFileOnUploadSuccess(file.GetName(), fileId, dstDir.GetPath(), fileParts, totalSize)
}
