use archive_reader::Archive;

use std::path::Path;
use archive_reader::error::Result;
use super::BookInfo;

fn valid_page(file: &String) -> bool {
    return file.to_lowercase().ends_with("jpg")
        || file.to_lowercase().ends_with("jpeg") 
        || file.to_lowercase().ends_with("png")
        || file.to_lowercase().ends_with("webp");
}

fn get_pages(list: &Vec<String>) -> usize {
    return list.into_iter().filter(|file: &&std::string::String| valid_page(*file)).count();
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
    let _excluded_files = file_names.clone().into_iter()
        .filter(|file| !valid_page(file))
        .for_each(|file| println!("excluded={:?}", file));

    return Ok(BookInfo {
        pages: get_pages(&file_names)
    });
}
