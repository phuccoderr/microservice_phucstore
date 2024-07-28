package com.phuc.categoryservice.exceptions;


import com.phuc.categoryservice.response.ResponseError;
import com.phuc.categoryservice.response.ResponseObject;
import io.jsonwebtoken.ExpiredJwtException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpStatus;
import org.springframework.http.HttpStatusCode;
import org.springframework.http.ResponseEntity;
import org.springframework.security.authentication.InsufficientAuthenticationException;
import org.springframework.validation.FieldError;
import org.springframework.web.bind.MethodArgumentNotValidException;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.bind.annotation.ResponseStatus;
import org.springframework.web.bind.annotation.RestControllerAdvice;
import org.springframework.web.context.request.WebRequest;
import org.springframework.web.servlet.mvc.method.annotation.ResponseEntityExceptionHandler;

import java.util.List;

@RestControllerAdvice
public class GlobalExceptionHandler extends ResponseEntityExceptionHandler {

    private static final Logger LOGGER = LoggerFactory.getLogger(GlobalExceptionHandler.class);

    @ExceptionHandler(Exception.class)
    public ResponseEntity<ResponseError> handleGeneralException(Exception ex) {

        LOGGER.error(ex.getMessage(),ex);

        return ResponseEntity.internalServerError().body(
            ResponseError.builder()
                    .statusCode(HttpStatus.INTERNAL_SERVER_ERROR.value())
                    .error(HttpStatus.INTERNAL_SERVER_ERROR.getReasonPhrase())
                    .message(List.of(ex.getMessage()))
                    .build()
        );
    }

    @ExceptionHandler(DataAlreadyExistsException.class)
    public ResponseEntity<ResponseError> handleDataAlreadyExistsException(Exception ex) {

        LOGGER.error(ex.getMessage(),ex);

        return ResponseEntity.unprocessableEntity().body(
                ResponseError.builder()
                        .statusCode(HttpStatus.UNPROCESSABLE_ENTITY.value())
                        .error(HttpStatus.UNPROCESSABLE_ENTITY.getReasonPhrase())
                        .message(List.of(ex.getMessage()))
                        .build()
        );
    }

    @ExceptionHandler(DataNotFoundException.class)
    public ResponseEntity<ResponseError> handleDataNotFoundException(Exception ex) {

        LOGGER.error(ex.getMessage(),ex);

        return ResponseEntity.unprocessableEntity().body(
                ResponseError.builder()
                        .statusCode(HttpStatus.UNPROCESSABLE_ENTITY.value())
                        .error(HttpStatus.UNPROCESSABLE_ENTITY.getReasonPhrase())
                        .message(List.of(ex.getMessage()))
                        .build()
        );
    }

    @ExceptionHandler(DataDuplicatedException.class)
    public ResponseEntity<ResponseError> handleDataDuplicatedException(Exception ex) {

        LOGGER.error(ex.getMessage(),ex);

        return ResponseEntity.unprocessableEntity().body(
                ResponseError.builder()
                        .statusCode(HttpStatus.UNPROCESSABLE_ENTITY.value())
                        .error(HttpStatus.UNPROCESSABLE_ENTITY.getReasonPhrase())
                        .message(List.of(ex.getMessage()))
                        .build()
        );
    }

    @ExceptionHandler(ParamValidateException.class)
    public ResponseEntity<ResponseError> handleParamValidateException(Exception ex) {

        LOGGER.error(ex.getMessage(),ex);

        return ResponseEntity.badRequest().body(
                ResponseError.builder()
                        .statusCode(HttpStatus.BAD_REQUEST.value())
                        .error(HttpStatus.BAD_REQUEST.getReasonPhrase())
                        .message(List.of(ex.getMessage()))
                        .build()
        );
    }

    @ExceptionHandler(DataHasChildrenException.class)
    public ResponseEntity<ResponseError> handleDataHasChildrenException(Exception ex) {

        LOGGER.error(ex.getMessage(),ex);

        return ResponseEntity.unprocessableEntity().body(
                ResponseError.builder()
                        .statusCode(HttpStatus.UNPROCESSABLE_ENTITY.value())
                        .error(HttpStatus.UNPROCESSABLE_ENTITY.getReasonPhrase())
                        .message(List.of(ex.getMessage()))
                        .build()
        );
    }

    @Override
    protected ResponseEntity<Object> handleMethodArgumentNotValid(
            MethodArgumentNotValidException ex,
            HttpHeaders headers,
            HttpStatusCode status,
            WebRequest request) {

        LOGGER.error(ex.getMessage(),ex);

        ResponseError responseError = ResponseError.builder()
                .error(HttpStatus.BAD_REQUEST.getReasonPhrase())
                .statusCode(HttpStatus.BAD_REQUEST.value()).build();

        List<FieldError> fieldError = ex.getBindingResult().getFieldErrors();
        fieldError.forEach( field -> responseError.addMessage(field.getDefaultMessage()));

        return new ResponseEntity<>(responseError, HttpStatus.BAD_REQUEST);
    }


}
