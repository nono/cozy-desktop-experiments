require "kemal"

macro my_render(filename)
  render "src/web/views/#{{{filename}}}.ecr", "src/web/views/layout.ecr"
end

module Web
  module Dev
    def self.setup_config
      config = Kemal.config
      config.host_binding = "127.0.0.1"
      config.port = 5177
      config.public_folder = "./src/web/public"
      config.setup
      config
    end

    def self.setup_router
      get "/" do |env|
        env.redirect "/generate"
      end

      get "/generate" do
        my_render "generate"
      end

      error 404 do
        render_404
      end
    end

    # XXX We can't use Kemal.run {|cfg| ... } as yield config is called after
    # config.setup, which means that options like public_folder are ignored.
    def self.run
      config = setup_config
      setup_router
      server = config.server ||= HTTP::Server.new(config.handlers)
      server.bind_tcp(config.host_binding, config.port)
      log "Listening on http://#{config.host_binding}:#{config.port}/"
      server.listen
    end
  end
end
