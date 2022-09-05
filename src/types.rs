use crate::types::response::DefaultResponse;
use actix_web::body::BoxBody;
use actix_web::http::StatusCode;
use actix_web::{HttpRequest, HttpResponse, Responder};
use serde::Serialize;
use std::fmt::{Debug, Display, Formatter};
use std::result;

pub mod models;
pub mod request;
pub mod response;

pub type Result<T> = result::Result<T, Error>;
pub struct Response<T>(Result<Option<T>>);

#[derive(Clone, Debug)]
pub enum Error {
    NotFound,
}

impl Display for Error {
    fn fmt(&self, formatter: &mut Formatter<'_>) -> std::fmt::Result {
        match self {
            Error::NotFound => formatter.write_str("these are not the droids you're looking for"),
        }
    }
}

impl Error {
    fn status_code(&self) -> StatusCode {
        match self {
            Error::NotFound => StatusCode::NOT_FOUND,
        }
    }
}

impl<T> Responder for Response<T>
where
    T: Serialize,
{
    type Body = BoxBody;

    fn respond_to(self, _request: &HttpRequest) -> HttpResponse<Self::Body> {
        let error = &self.0.as_ref().err().unwrap_or(&Error::NotFound).clone();
        let body = self.0.unwrap();

        return if body.is_none() {
            HttpResponse::build(error.status_code()).json(DefaultResponse {
                message: error.to_string(),
            })
        } else {
            HttpResponse::Ok().json(body.unwrap())
        };
    }
}
