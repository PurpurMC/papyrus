use crate::Config;
use actix_files::Files;
use actix_web::middleware::{Logger, NormalizePath};
use actix_web::web::Data;
use actix_web::{web, App, HttpServer};
use sqlx::SqlitePool;
use std::io::Result;

mod middleware;
mod router;

pub async fn spin(config: Config, pool: SqlitePool) -> Result<()> {
    let cloned_config = config.clone();
    HttpServer::new(move || {
        let config = cloned_config.clone();

        let mut app = App::new()
            .wrap(Logger::default())
            .wrap(NormalizePath::trim())
            .app_data(Data::new(pool.clone()))
            .app_data(Data::new(config.clone()))
            .service(web::scope("/v2").configure(router::setup));

        if config.docs.enabled {
            app = app.service(Files::new("/", config.docs.path).index_file("index.html"));
        }

        return app;
    })
    .bind(config.host)?
    .run()
    .await
}
