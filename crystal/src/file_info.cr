struct Crystal::System::FileInfo < ::File::Info
  # Return the inode number (on UNIX, not windows)
  def ino
    @stat.st_ino
  end
end
