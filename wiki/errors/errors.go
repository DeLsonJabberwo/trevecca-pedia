package errors

import (
	"errors"
	"net/http"
)

type WikiError struct {
	Code    int
	Type    string
	Details string
	err     error
}

func (w WikiError) Error() string {
	return w.Details
}

func (w WikiError) Unwrap() error {
	return w.err
}

// Error Types
const (
	pageNotFound     	= "PageNotFound"
	revisonNotFound  	= "RevisionNotFound"
	snapshotNotFound 	= "SnapshotNotFound"
	pageDeleted      	= "PageDeleted"
	revisionDeleted  	= "RevisionDeleted"
	snapshotDeleted  	= "SnapshotDeleted"
	invalidId        	= "InvalidId"
	revisionConflict 	= "RevisionConflict"
	internalErr			= "InternalServerError"
	databaseErr         = "DatabaseError"
	filesystemErr       = "FilesystemError"
	dbfsErr				= "DatabaseFilesystemError"
)

// Error Constructors
func PageNotFound() WikiError {
	return WikiError{
		http.StatusNotFound,
		pageNotFound,
		"page not found",
		nil,
	}
}

func RevisionNotFound() WikiError {
	return WikiError{
		http.StatusNotFound,
		revisonNotFound,
		"revision not found",
		nil,
	}
}

func SnapshotNotFound() WikiError {
	return WikiError{
		http.StatusNotFound,
		snapshotNotFound,
		"snapshot not found",
		nil,
	}
}

func PageDeleted() WikiError {
	return WikiError{
		http.StatusNotFound,
		pageDeleted,
		"page not found",
		nil,
	}
}

func RevisionDeleted() WikiError {
	return WikiError{
		http.StatusNotFound,
		revisionDeleted,
		"revision not found",
		nil,
	}
}

func SnapshotDeleted() WikiError {
	return WikiError{
		http.StatusNotFound,
		snapshotDeleted,
		"snapshot not found",
		nil,
	}
}

func InvalidID(err error) WikiError {
	return WikiError{
		http.StatusBadRequest,
		invalidId,
		"invalid id",
		err,
	}
}

func RevisionConflict(err error) WikiError {
	return WikiError{
		http.StatusConflict,
		revisionConflict,
		"revision conflict",
		err,
	}
}

func InternalError(err error) WikiError {
	return WikiError{
		http.StatusInternalServerError,
		internalErr,
		"internal server error",
		err,
	}
}

func DatabaseError(err error) WikiError {
	return WikiError{
		http.StatusInternalServerError,
		databaseErr,
		"database error",
		err,
	}
}

func FilesystemError(err error) WikiError {
	return WikiError{
		http.StatusInternalServerError,
		filesystemErr,
		"filesystem error",
		err,
	}
}

func DatabaseFilesystemError(err error) WikiError {
	return WikiError{
		http.StatusInternalServerError,
		dbfsErr,
		"database/filesystem error",
		err,
	}
}


// Utilities
func IsWikiError(err error) bool {
	if err == nil {
		return false
	}
	var we WikiError
	return errors.As(err, &we)
}

func AsWikiError(err error) (WikiError, bool) {
	if err == nil {
		return WikiError{}, false
	}
	var we WikiError
	if !errors.As(err, &we) {
		return WikiError{}, false
	}
	return we, true
}

func HasType(err error, errType string) bool {
	we, ok := AsWikiError(err)
	if !ok {
		return false
	}
	return we.Type == errType
}

func IsNotFound(err error) bool {
	return HasType(err, pageNotFound) ||
		HasType(err, revisonNotFound) ||
		HasType(err, snapshotNotFound)
}

func IsDeleted(err error) bool {
	return HasType(err, pageDeleted) ||
		HasType(err, revisionDeleted) ||
		HasType(err, snapshotDeleted)
}
