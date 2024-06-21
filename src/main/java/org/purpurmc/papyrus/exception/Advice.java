package org.purpurmc.papyrus.exception;

import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.ControllerAdvice;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.bind.annotation.ResponseBody;
import org.springframework.web.bind.annotation.ResponseStatus;
import org.springframework.web.servlet.NoHandlerFoundException;

@ControllerAdvice
public class Advice {
    @ExceptionHandler(BuildAlreadyExists.class)
    @ResponseBody
    @ResponseStatus(HttpStatus.CONFLICT)
    public ErrorResponse buildAlreadyExists() {
        return new ErrorResponse("build already exists");
    }

    @ExceptionHandler(BuildNotFound.class)
    @ResponseBody
    @ResponseStatus(HttpStatus.NOT_FOUND)
    public ErrorResponse buildNotFound() {
        return new ErrorResponse("build not found");
    }

    @ExceptionHandler(FileDownloadError.class)
    @ResponseBody
    @ResponseStatus(HttpStatus.INTERNAL_SERVER_ERROR)
    public ErrorResponse fileDownload() {
        return new ErrorResponse("couldn't access file");
    }

    @ExceptionHandler(FileUploadError.class)
    @ResponseBody
    @ResponseStatus(HttpStatus.INTERNAL_SERVER_ERROR)
    public ErrorResponse fileUpload() {
        return new ErrorResponse("couldn't upload file");
    }

    @ExceptionHandler(InvalidAuthToken.class)
    @ResponseBody
    @ResponseStatus(HttpStatus.FORBIDDEN)
    public ErrorResponse invalidAuthToken() {
        return new ErrorResponse("invalid auth token");
    }

    @ExceptionHandler(InvalidStateKey.class)
    @ResponseBody
    @ResponseStatus(HttpStatus.BAD_REQUEST)
    public ErrorResponse invalidStateKey() {
        return new ErrorResponse("invalid state key");
    }

    @ExceptionHandler(ProjectNotFound.class)
    @ResponseBody
    @ResponseStatus(HttpStatus.NOT_FOUND)
    public ErrorResponse projectNotFound() {
        return new ErrorResponse("project not found");
    }

    @ExceptionHandler(VersionNotFound.class)
    @ResponseBody
    @ResponseStatus(HttpStatus.NOT_FOUND)
    public ErrorResponse versionNotFound() {
        return new ErrorResponse("version not found");
    }

    @ExceptionHandler(NoHandlerFoundException.class)
    @ResponseBody
    @ResponseStatus(HttpStatus.NOT_FOUND)
    public ErrorResponse noHandlerFound() {
        return new ErrorResponse("endpoint not found");
    }

    private record ErrorResponse(String error) {
    }
}
