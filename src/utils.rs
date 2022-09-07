use crate::types::{Error, Result};
use actix_web::HttpResponse;
use serde::Serialize;

pub fn response<T>(body: T) -> Result<HttpResponse>
where
    T: Serialize,
{
    Ok(HttpResponse::Ok().json(body))
}

pub fn verify<T>(option: Option<T>) -> Result<T> {
    match option {
        Some(value) => Ok(value),
        None => Err(Error::NotFound),
    }
}
