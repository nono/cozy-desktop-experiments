# TODO: try [spider-gazelle](https://spider-gazelle.net/#/README) instead of kemal
require "kemal"
require "../simulator/*"

macro my_render(filename)
  render "src/web/views/#{{{filename}}}.ecr", "src/web/views/layout.ecr"
end

module Web
  module Dev
    # XXX We can't use Kemal.run {|cfg| ... } as yield config is called after
    # config.setup, which means that options like public_folder are ignored.
    def self.run
      config = setup_config
      setup_router
      server = config.server ||= HTTP::Server.new(config.handlers)
      server.bind_tcp(config.host_binding, config.port)
      config.running = true
      log "Listening on http://#{config.host_binding}:#{config.port}/"
      server.listen
    end

    def self.setup_config
      config = Kemal.config
      config.host_binding = "127.0.0.1"
      config.port = 5177
      config.public_folder = "./src/web/public"
      config.setup
      config
    end

    def self.setup_router
      get "/" do
        my_render "index"
      end

      get "/generate" do
        my_render "generate"
      end

      post "/generate" do |env|
        env.response.content_type = "application/json"
        nb_clients = env.params.json["nb_clients"].as(Int64)
        scenario = Simulator::Generate.new(nb_clients).run
        pp scenario.to_json
      end

      post "/play" do |env|
        env.response.content_type = "application/json"
        scenario = Simulator::Generate.new(1).run # TODO: load the scenario from the request
        play = Simulator::Play.new(scenario)
        play.run
        errors = play.check
        pp errors.to_json
      end

      error 404 do
        render_404
      end
    end
  end
end
