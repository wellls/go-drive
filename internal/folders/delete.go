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

	err = deleteFolderContent(h.db, int64(id))
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

func deleteFolderContent(db *sql.DB, folderID int64) error {
	err := deleteFiles(db, folderID)
	if err != nil {
		return err
	}

	return deleteSubFolders(db, folderID)
}

func deleteSubFolders(db *sql.DB, folderID int64) error {

	subfolders, err := getSubFolders(db, folderID)
	if err != nil {
		return err
	}

	removedFolders := make([]Folder, 0, len(subfolders))
	for _, sf := range subfolders {
		err := Delete(db, sf.ID)
		if err != nil {
			break
		}

		err = deleteFolderContent(db, sf.ID)
		if err != nil {
			Update(db, sf.ID, &sf)
			break
		}

		removedFolders = append(removedFolders, sf)
	}

	if len(subfolders) != len(removedFolders) {
		for _, rf := range removedFolders {
			Update(db, rf.ID, &rf)
		}
	}

	return nil

}

func deleteFiles(db *sql.DB, folderID int64) error {
	f, err := files.List(db, folderID)
	if err != nil {
		return err
	}

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
