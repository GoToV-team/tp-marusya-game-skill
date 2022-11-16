package prepare

import (
    "github.com/pkg/errors"
    "os"
    "path/filepath"
    "time"
)

const logName = "log"

func OpenLogDir(dir string) (*os.File, error) {
    if _, err := os.Stat(dir); err != nil {
        if !os.IsNotExist(err) {
            return nil, errors.Wrap(err, "error when try check log dir: ")
        }

        if err = os.MkdirAll(filepath.Dir(dir), 0755); err != nil {
            return nil, errors.Wrap(err, "error when try add log dir: ")
        }
    }

    t := time.Now().UTC()
    timeString := t.Format(time.RFC3339)
    fileName := timeString + "-" + logName + ".log"

    file, err := os.OpenFile(
        dir+"/"+fileName,
        os.O_CREATE|os.O_APPEND|os.O_WRONLY,
        0644,
    )

    if err != nil {
        return nil, errors.Wrap(err, "error when try open log file: ")
    }

    return file, nil
}
