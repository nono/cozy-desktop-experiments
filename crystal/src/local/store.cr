require "db"
require "sqlite3"

module Local
  class Store
    property scan_counter
    property db : DB::Database

    def initialize
      @scan_counter = 0
      # TODO: use a persistent database
      @db = DB.open "sqlite3://%3Amemory%3A"
    end

    # setup creates the SQL table for the local files and directories.
    #
    # The inode is not used as the primary key, as Linux can reuse an inode for
    # a new file, and we may want to keep the same ID for a file saved with the
    # atomic technic (ie writing to a new temprorary file and moving this file
    # over the existing file).
    #
    # AUTOINCREMENT is used for the primary key to tells sqlite to not reuse
    # IDs from deleted rows. See https://www.sqlite.org/autoinc.html
    #
    # The path is not kept in the database, only the identifier of the parent
    # directory. It may change later, but it looks simpler this way (as it
    # avoids to manipulate many rows when a directory is renamed/moved), and I
    # would try to see if there is a real drawback to not have the path.
    def setup
      @db.exec "CREATE TABLE IF NOT EXISTS local (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        parent_id INTEGER,
        inode INTEGER NOT NULL,
        type TEXT NOT NULL CHECK(type IN ('dir', 'file')),
        name TEXT NOT NULL
      )"
    end

    def close
      @db.close
    end
  end
end
