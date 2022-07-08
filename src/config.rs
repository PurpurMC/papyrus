use std::fs::File;
use std::net::SocketAddr;
use std::path::Path;
use serde::{Deserialize, Serialize};
use crate::utils::Error;

#[derive(Serialize, Deserialize)]
pub struct Config {
    pub host: SocketAddr,
    pub database: String,
}

impl Config {
    pub fn default() -> Config {
        Config {
            host: SocketAddr::new("0.0.0.0".parse().unwrap(), 8080),
            database: "/srv/papyrus".into()
        }
    }

    pub fn load() -> Result<Config, Error> {
        let config_path = Path::new("config.json"); // todo

        if !config_path.exists() {
            let config = Config::default();
            let file = File::create(config_path)?;
            serde_json::to_writer_pretty(file, &config)?;
            Ok(config)
        } else {
            let file = File::open(config_path)?;
            let config: Config = serde_json::from_reader(file)?;
            Ok(config)
        }
    }
}
