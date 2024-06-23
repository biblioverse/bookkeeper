use super::BookInfo;
use archive_reader::error::Result;
use archive_reader::Archive;
use std::ffi::OsStr;
use std::path::Path;

fn valid_page(file: &str) -> bool {
    file.to_lowercase().ends_with("jpg")
        || file.to_lowercase().ends_with("jpeg")
        || file.to_lowercase().ends_with("png")
        || file.to_lowercase().ends_with("webp")
}

fn get_pages(list: &[String]) -> usize {
    list.iter()
        .filter(|file: &&std::string::String| valid_page(file))
        .count()
}

pub fn get_book_info_cba(file: &Path) -> Result<BookInfo, Box<dyn std::error::Error>> {
    let mut archive = Archive::open(file);
    let file_names = archive
        .block_size(1024 * 1024)
        .list_file_names()?
        .collect::<Result<Vec<_>>>()?;
    //let mut content = vec![];
    //let _ = archive.read_file(&file_names[0], &mut content)?;
    //println!("content={:?}", &file_names[0]);

    // Find metadata file
    file_names
        .clone()
        .into_iter()
        .filter(|file| !valid_page(file))
        .for_each(|file| println!("excluded={:?}", file));

    Ok(BookInfo {
        title: format!("{:?}", file.file_stem().unwrap_or(OsStr::new("stuff"))),
        pages: get_pages(&file_names),
        authors: None,
        publisher: None,
        published_date: None,
        keywords: None,
    })
}
