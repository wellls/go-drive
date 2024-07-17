package folders

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/wellls/go-drive/internal/files"
)

func (h *handler) Delete(rw http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	err = deleteFiles(h.db, int64(id))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	err = Delete(h.db, int64(id))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Add("Content-Type", "application/json")
}

func deleteFiles(db *sql.DB, folderID int64) error {
	f, err := files.List(db, folderID)
	if err != nil {
		return err
	}

	// TODO: List Folders

	removedFiles := make([]files.File, 0, len(f))
	for _, file := range f {
		file.Deleted = true
		err = files.Update(db, file.ID, &file)
		if err != nil {
			break
		}

		removedFiles = append(removedFiles, file)
	}

	if len(f) != len(removedFiles) {
		for _, file := range removedFiles {
			file.Deleted = false
			err = files.Update(db, file.ID, &file)
		}

		return err
	}

	return nil
}

func Delete(db *sql.DB, id int64) error {
	stmt := `update "folders" set "modified_at"=$1, "deleted"=true where id=$2`
	_, err := db.Exec(stmt, time.Now(), id)

	return err
}
