use crate::Config;
use actix_web::body::BoxBody;
use actix_web::dev::{forward_ready, Service, ServiceRequest, ServiceResponse, Transform};
use actix_web::web::Data;
use actix_web::{Error, HttpResponse};
use futures_util::future::{ok, LocalBoxFuture, Ready};
use serde_json::json;

pub struct Authentication;

impl<S> Transform<S, ServiceRequest> for Authentication
where
    S: Service<ServiceRequest, Response = ServiceResponse<BoxBody>, Error = Error>,
    S::Future: 'static,
{
    type Response = ServiceResponse<BoxBody>;
    type Error = Error;
    type Transform = AuthenticationMiddleware<S>;
    type InitError = ();
    type Future = Ready<Result<Self::Transform, Self::InitError>>;

    fn new_transform(&self, service: S) -> Self::Future {
        ok(AuthenticationMiddleware { service })
    }
}

pub struct AuthenticationMiddleware<S> {
    service: S,
}

impl<S> Service<ServiceRequest> for AuthenticationMiddleware<S>
where
    S: Service<ServiceRequest, Response = ServiceResponse<BoxBody>, Error = Error>,
    S::Future: 'static,
{
    type Response = ServiceResponse<BoxBody>;
    type Error = Error;
    type Future = LocalBoxFuture<'static, Result<Self::Response, Self::Error>>;

    forward_ready!(service);

    fn call(&self, request: ServiceRequest) -> Self::Future {
        let passed = loop {
            let config = match request.app_data::<Data<Config>>() {
                Some(config) => config,
                None => break false,
            };

            let auth = match request.headers().get("Authorization") {
                Some(auth) => auth,
                None => break false,
            };

            let auth_parts = match auth.to_str() {
                Ok(auth_parts) => auth_parts.split(' ').collect::<Vec<&str>>(),
                Err(_) => break false,
            };

            if auth_parts.len() != 2 || auth_parts[0] != "Token" {
                break false;
            }

            let key = auth_parts[1];
            if config.auth_key != key {
                break false;
            }

            break true;
        };

        if passed {
            let future = self.service.call(request);
            Box::pin(async move {
                let res = future.await?;
                Ok(res)
            })
        } else {
            let (request, _pl) = request.into_parts();
            let response = HttpResponse::Found().json(json!({ "error": "Unauthorized" }));

            Box::pin(async move { Ok(ServiceResponse::new(request, response)) })
        }
    }
}
