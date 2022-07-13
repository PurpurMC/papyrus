use crate::config::Config;
use actix_web::middleware::Logger;
use actix_web::web::Data;
use actix_web::{web, App, HttpServer};
use dotenv::dotenv;
use env_logger::Env;
use sqlx::sqlite::SqliteConnectOptions;
use sqlx::SqlitePool;
use std::str::FromStr;
use std::{fs, io};

mod config;
mod models;
mod router;
mod utils;

#[actix_web::main]
async fn main() -> io::Result<()> {
    dotenv().ok();
    env_logger::init_from_env(Env::new().default_filter_or("info"));

    let config = Config::load().map_err(|e| io::Error::new(io::ErrorKind::Other, e))?;

    let _ = fs::create_dir(&config.database); // todo: im too lazy to handle this error
    let pool = SqlitePool::connect_with(
        SqliteConnectOptions::from_str(&*format!("sqlite://{}/database.db", &config.database))
            .map_err(|e| io::Error::new(io::ErrorKind::Other, e))?
            .create_if_missing(true),
    )
    .await
    .map_err(|e| io::Error::new(io::ErrorKind::Other, e))?;

    sqlx::migrate!()
        .run(&pool)
        .await
        .map_err(|e| io::Error::new(io::ErrorKind::Other, e))?;

    HttpServer::new(move || {
        App::new()
            .wrap(Logger::default())
            .app_data(Data::new(pool.clone()))
            .app_data(Data::new(Config::load().unwrap()))
            .service(web::scope("/v2").configure(router::setup))
    })
    .bind(config.host)?
    .run()
    .await
}
