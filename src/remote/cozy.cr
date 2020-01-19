require "http/client"
require "uri"

module Remote
  class Cozy
    ROOT_DIR_ID  = "io.cozy.files.root-dir"
    TRASH_DIR_ID = "io.cozy.files.trash-dir"

    def initialize(url : String, access_token : String)
      # TODO: set dns/connect/read/write timeouts
      @client = HTTP::Client.new url
      # TODO: allow to refresh the access_token
      @client.before_request do |request|
        request.headers["Authorization"] = "Bearer #{access_token}"
      end
    end

    def createDirectory(dir_id, name, date)
      headers = HTTP::Headers{"Date" => date}
      path = "/files/#{dir_id}?Type=directory&Name=#{URI.encode name}"
      @client.post path, headers
      # TODO: raise an exception if the response is not a success
      # TODO: parse the response body as JSON
    end
  end
end
