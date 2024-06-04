use super::utils::error::IllegalArgumentError;
use super::utils::fileutils::{is_pdf, is_zip};
use cb::get_book_info_cba;
use pdf::get_book_info_pdf;
use serde::{Deserialize, Serialize};
use std::error::Error;
use std::path::Path;

pub mod cb;
pub mod pdf;

#[derive(Serialize, Deserialize)]
pub struct BookInfo {
    title: String,
    pages: usize,
    authors: Option<Vec<String>>,
    publisher: Option<String>,
    published_date: Option<String>,
    keywords: Option<Vec<String>>,
}

pub fn get_book_info(file: &Path) -> Result<BookInfo, Box<dyn Error>> {
    if is_zip(file) {
        return get_book_info_cba(file);
    }

    if is_pdf(file) {
        return get_book_info_pdf(file);
    }

    Err(Box::new(IllegalArgumentError::new(format!(
        "We don't know how to open this archive '{:?}'",
        file
    ))))
}
