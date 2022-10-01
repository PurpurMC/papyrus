use crate::types::response::DefaultResponse;
use crate::Config;
use actix_web::body::BoxBody;
use actix_web::dev::{forward_ready, Service, ServiceRequest, ServiceResponse, Transform};
use actix_web::web::Data;
use actix_web::{Error, HttpResponse};
use futures_util::future::{ok, LocalBoxFuture, Ready};

pub struct RequireKey;

impl<S> Transform<S, ServiceRequest> for RequireKey
where
    S: Service<ServiceRequest, Response = ServiceResponse<BoxBody>, Error = Error>,
    S::Future: 'static,
{
    type Response = ServiceResponse<BoxBody>;
    type Error = Error;
    type Transform = RequireKeyMiddleware<S>;
    type InitError = ();
    type Future = Ready<Result<Self::Transform, Self::InitError>>;

    fn new_transform(&self, service: S) -> Self::Future {
        ok(RequireKeyMiddleware { service })
    }
}

pub struct RequireKeyMiddleware<S> {
    service: S,
}

impl<S> Service<ServiceRequest> for RequireKeyMiddleware<S>
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

            let header = match request.headers().get("Authorization") {
                Some(auth) => auth,
                None => break false,
            };

            let header_parts = match header.to_str() {
                Ok(auth_parts) => auth_parts.split(' ').collect::<Vec<&str>>(),
                Err(_) => break false,
            };

            break header_parts.len() == 2
                && header_parts[0].eq("Token")
                && config.keys.contains(&header_parts[1].to_string())
        };

        if passed {
            let future = self.service.call(request);
            Box::pin(async move { Ok(future.await?) })
        } else {
            let (request, _payload) = request.into_parts();
            let response = HttpResponse::Unauthorized().json(DefaultResponse {
                message: "Invalid authorization key".into(),
            });

            Box::pin(async move { Ok(ServiceResponse::new(request, response)) })
        }
    }
}
