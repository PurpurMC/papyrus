use crate::server::middleware::RequireKey;
use actix_web::web;
use actix_web::web::ServiceConfig;

pub fn routes(config: &mut ServiceConfig) {
    config.service(web::scope("").wrap(RequireKey));
}
