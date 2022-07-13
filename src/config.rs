use crate::utils::Error;
use nanoid::nanoid;
use serde::{Deserialize, Serialize};
use std::fs::File;
use std::net::SocketAddr;
use std::path::Path;

#[derive(Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct Config {
    pub host: SocketAddr,
    pub database: String,
    pub auth_key: String,
}

impl Config {
    pub fn default() -> Config {
        Config {
            host: SocketAddr::new("0.0.0.0".parse().unwrap(), 8080),
            database: "/srv/papyrus".into(),
            auth_key: nanoid!(128),
        }
    }

    pub fn load() -> Result<Config, Error> {
        let config_path = Path::new(if cfg!(debug_assertions) {
            "config.json"
        } else {
            "/etc/papyrus.json"
        });

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
