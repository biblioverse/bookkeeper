use serde::{Deserialize, Serialize};
use sha1::{Digest, Sha1};
use std::{
    fs::{self, File},
    io::{BufReader, Read},
    path::Path,
};
use walkdir::{DirEntry, WalkDir};

use super::super::archives;
use super::super::archives::BookInfo;
use super::super::utils::fileutils::{is_pdf, is_zip};

#[derive(Serialize, Deserialize)]
struct Metadata {
    path: String,
    status: String,
    size: u64,
    hash: String,
    book: BookInfo,
}

#[derive(Serialize, Deserialize)]
struct ErrorMetadata {
    path: String,
    status: String,
    error: String,
}

fn is_valid_file(entry: &DirEntry) -> bool {
    if entry.file_type().is_dir() {
        return false;
    }

    return is_zip(entry.path()) || is_pdf(entry.path());
}

fn hash_file(path: &std::path::Path) -> Result<String, Box<dyn std::error::Error>> {
    if true {
        return Ok(String::new());
    }

    // Open the file
    let file = File::open(path)?;
    let mut reader = BufReader::new(file);

    // Create a Sha1 object
    let mut hasher = Sha1::new();

    // Read the file in chunks and update the hash
    let mut buffer = [0; 1024];
    loop {
        let bytes_read = reader.read(&mut buffer)?;
        if bytes_read == 0 {
            break;
        }
        hasher.update(&buffer[..bytes_read]);
    }

    // Finalize the hash
    let result = hasher.finalize();
    let hash_string = format!("{:x}", result);

    Ok(hash_string)
}

fn scan_book(full_scan_path: &Path, filepath: &Path) -> Result<String, Box<dyn std::error::Error>> {
    let book = archives::get_book_info(filepath)?;

    let metadata = fs::metadata(filepath)?;
    let computed = Metadata {
        path: format!("{:?}", filepath.strip_prefix(full_scan_path)?),
        status: "success".to_string(),
        hash: hash_file(filepath)?,
        size: metadata.len().to_owned(),
        book,
    };

    Ok(serde_json::to_string(&computed)?)
}

fn print_error(
    full_scan_path: &Path,
    filepath: &Path,
    error: Box<dyn std::error::Error>,
) -> String {
    let split_path = match filepath.strip_prefix(full_scan_path) {
        Ok(path) => format!("{:?}", path),
        Err(err) => format!("Could not serialize path {:?}", err),
    };

    let json = ErrorMetadata {
        path: split_path,
        status: "failed".to_string(),
        error: format!("{:?}", error),
    };

    match serde_json::to_string(&json) {
        Ok(r) => r,
        Err(err) => format!("Could not serialize {:?}", err),
    }
}

pub fn scan(scan_path: &std::path::PathBuf) -> Result<(), Box<dyn std::error::Error>> {
    let canonical_path = fs::canonicalize(scan_path)?;
    let full_scan_path = canonical_path.as_path();

    for entry in WalkDir::new(full_scan_path)
        .into_iter()
        .filter_map(Result::ok)
        .filter(is_valid_file)
    {
        let filepath = entry.path();

        match scan_book(full_scan_path, filepath) {
            Ok(line) => println!("{}", line),
            Err(error) => eprintln!("{}", print_error(full_scan_path, filepath, error)),
        }
    }

    Ok(())
}
