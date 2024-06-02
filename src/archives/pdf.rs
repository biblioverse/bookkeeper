use lopdf::Document;

use std::path::Path;
use archive_reader::error::Result;
use super::BookInfo;

pub fn get_book_info_pdf(file: &Path) -> Result<BookInfo, Box<dyn std::error::Error>> {
    let doc = Document::load(file)?;
    let pages = doc.get_pages();

    return Ok(BookInfo {
        pages: pages.len()
    })
}
