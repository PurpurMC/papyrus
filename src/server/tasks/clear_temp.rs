use crate::types::models::TempFile;
use crate::types::Result;
use crate::{Config, SqlitePool};
use async_trait::async_trait;
use std::fs;
use std::path::Path;
use std::time::Duration;

pub struct ClearTemp {
    pub config: Config,
    pub pool: SqlitePool,
}

#[async_trait]
impl super::Task for ClearTemp {
    fn duration(&self) -> Duration {
        Duration::from_secs(60 * 60)
    }

    async fn run(&self) -> Result<()> {
        let temp_files: Vec<TempFile> = sqlx::query_as!(
            TempFile,
            "SELECT * FROM temp_files WHERE created_at <= date('now', '-1 day')"
        )
        .fetch_all(&self.pool)
        .await?;

        for temp_file in temp_files {
            let path = format!("{}/temp/{}", &self.config.database, &temp_file.id);
            let path = Path::new(&path);
            if path.exists() {
                fs::remove_file(path)?;
            }

            TempFile::delete(&temp_file.id, &self.pool).await?;
        }

        Ok(())
    }
}
