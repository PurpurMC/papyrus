use crate::types::response::DefaultResponse;
use actix_web::http::StatusCode;
use actix_web::{HttpResponse, Responder, ResponseError};
use std::fmt::{Debug, Formatter};
use std::result;
use thiserror::Error;

pub mod models;
pub mod request;
pub mod response;

pub type Result<T> = result::Result<T, Error>;

#[derive(Error, Debug)]
pub enum Error {
    #[error("these are not the droids you're looking for")]
    NotFound,

    #[error(transparent)]
    DatabaseError(#[from] sqlx::Error),
}

impl ResponseError for Error {
    fn status_code(&self) -> StatusCode {
        match self {
            Error::NotFound => StatusCode::NOT_FOUND,
            Error::DatabaseError(_) => StatusCode::INTERNAL_SERVER_ERROR,
        }
    }

    fn error_response(&self) -> HttpResponse {
        HttpResponse::build(self.status_code()).json(DefaultResponse {
            message: self.to_string(),
        })
    }
}
