use crate::config::Config;
use dotenv::dotenv;
use env_logger::Env;
use sqlx::sqlite::SqliteConnectOptions;
use sqlx::SqlitePool;
use std::str::FromStr;
use std::{fs, io};

mod config;
mod server;
mod types;

#[actix_web::main]
async fn main() -> io::Result<()> {
    dotenv().ok();
    env_logger::init_from_env(Env::new().default_filter_or("info"));

    let (existed, config) = Config::load()?;
    if !existed {
        println!("The config file has been created, please edit it and restart the server");
        return Ok(());
    }

    fs::create_dir_all(format!("{}/files", &config.database))?;

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

    server::spin(config, pool).await
}
