use crate::server::tasks::{ClearTemp, Task};
use crate::Config;
use actix_files::Files;
use actix_web::middleware::{Logger, NormalizePath};
use actix_web::rt::time;
use actix_web::web::Data;
use actix_web::{rt, web, App, HttpServer};
use sqlx::SqlitePool;
use std::io::Result;

mod middleware;
mod router;
mod tasks;

pub async fn spin(config: Config, pool: SqlitePool) -> Result<()> {
    let tasks: Vec<Box<dyn Task>> = vec![Box::new(ClearTemp {
        config: config.clone(),
        pool: pool.clone(),
    })];

    for task in tasks {
        rt::spawn(async move {
            let mut interval = time::interval(task.duration());
            loop {
                interval.tick().await;
                let _ = task.run().await;
            }
        });
    }

    let server_config = config.clone();
    HttpServer::new(move || {
        let config = server_config.clone();

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
