use crate::types::Result;
use crate::SqlitePool;
use chrono::NaiveDateTime;
use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize)]
pub struct Project {
    pub id: String,
    pub name: String,
    pub created_at: NaiveDateTime,
}

pub struct ProjectSearch {
    pub name: String,
}

impl super::Model<ProjectSearch> for Project {
    fn all(pool: &SqlitePool) -> Result<Vec<Box<Self>>> {
        todo!()
    }

    fn find_by(params: ProjectSearch, pool: &SqlitePool) -> Result<Option<Box<Self>>> {
        todo!()
    }
}

#[derive(Serialize, Deserialize)]
pub struct Version {
    pub id: String,
    pub name: String,
    pub project_id: String,
    pub created_at: NaiveDateTime,
}

pub struct VersionSearch {
    pub name: String,
}

impl super::Model<VersionSearch> for Version {
    fn all(pool: &SqlitePool) -> Result<Vec<Box<Self>>> {
        todo!()
    }

    fn find_by(params: VersionSearch, pool: &SqlitePool) -> Result<Option<Box<Self>>> {
        todo!()
    }
}

#[derive(Serialize, Deserialize)]
pub struct Build {
    pub id: String,
    pub name: String,
    pub version_id: String,
    pub result: String,
    pub duration: i64,
    pub timestamp: i64,
    pub created_at: NaiveDateTime,
}

pub struct BuildSearch {
    pub name: String,
}

impl super::Model<BuildSearch> for Build {
    fn all(pool: &SqlitePool) -> Result<Vec<Box<Self>>> {
        todo!()
    }

    fn find_by(params: BuildSearch, pool: &SqlitePool) -> Result<Option<Box<Self>>> {
        todo!()
    }
}

#[derive(Serialize, Deserialize)]
pub struct File {
    pub id: String,
    pub build_id: Option<String>,
    pub hash: String,
    pub extension: String,
    pub created_at: NaiveDateTime,
}

pub struct FileSearch {
    pub name: String,
}

impl super::Model<FileSearch> for File {
    fn all(pool: &SqlitePool) -> Result<Vec<Box<Self>>> {
        todo!()
    }

    fn find_by(params: FileSearch, pool: &SqlitePool) -> Result<Option<Box<Self>>> {
        todo!()
    }
}

#[derive(Serialize, Deserialize)]
pub struct Commit {
    pub id: String,
    pub build_id: String,
    pub author: String,
    pub email: String,
    pub description: String,
    pub hash: String,
    pub timestamp: i64,
}

pub struct CommitSearch {
    pub name: String,
}

impl super::Model<CommitSearch> for Commit {
    fn all(pool: &SqlitePool) -> Result<Vec<Box<Self>>> {
        todo!()
    }

    fn find_by(params: CommitSearch, pool: &SqlitePool) -> Result<Option<Box<Self>>> {
        todo!()
    }
}
