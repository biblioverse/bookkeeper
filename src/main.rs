use clap::{Parser, Subcommand};

mod commands;
mod archives;
mod utils;

use commands::scan;

#[derive(Parser)]
#[command(author, version, about, long_about = None)]
#[command(propagate_version = true)]
struct Cli {
    #[command(subcommand)]
    command: Option<Commands>,
}

#[derive(Subcommand)]
enum Commands {
    /// Scan a folder recursively for books
    Scan { path: std::path::PathBuf },
}

fn main() {
    let cli = Cli::parse();

    // You can check for the existence of subcommands, and if found use their
    // matches just as you would the top level cmd
    match &cli.command {
        Some(Commands::Scan { path }) => {
            println!("Scanning: {:?}", path);
            let scan_result = scan::scan(path);

            match scan_result {
                Ok(_) => { println!("Scanned all files"); }
                Err(error) => { println!("Oh noes: {}", error); }
            }
        }
        None => {
            println!("Default subcommand");
        }
    }
}