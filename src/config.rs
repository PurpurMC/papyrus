use nanoid::nanoid;
use serde::{Deserialize, Serialize};
use std::fs::File;
use std::io::Error;
use std::net::SocketAddr;
use std::path::Path;

#[derive(Serialize, Deserialize, Clone)]
#[serde(rename_all = "camelCase")]
pub struct Config {
    pub host: SocketAddr,
    pub database: String,
    pub auth_key: String, // todo: better auth keys system
    pub docs: Docs,
}

#[derive(Serialize, Deserialize, Clone)]
#[serde(rename_all = "camelCase")]
pub struct Docs {
    pub enabled: bool,
    pub path: String,
}

impl Config {
    pub fn default() -> Config {
        Config {
            host: SocketAddr::new("0.0.0.0".parse().unwrap(), 8080),
            database: "/srv/papyrus".into(),
            auth_key: nanoid!(128),
            docs: Docs {
                enabled: false,
                path: "/var/www/papyrus".into(),
            },
        }
    }

    pub fn load() -> Result<(bool, Config), Error> {
        let config_path = Path::new(if cfg!(debug_assertions) {
            "config.json"
        } else {
            "/etc/papyrus.json"
        });

        if !config_path.exists() {
            let config = Config::default();
            serde_json::to_writer_pretty(File::create(config_path)?, &config)?;
            Ok((true, config))
        } else {
            let config: Config = serde_json::from_reader(File::open(config_path)?)?;
            Ok((false, config))
        }
    }
}
