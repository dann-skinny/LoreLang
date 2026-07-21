require "sinatra/base"
require "json"
require "net/http"
require "uri"

require_relative "tom_nook_agent"

class LLMClient
  OLLAMA_URL = URI("http://localhost:11434/api/chat")

  def initialize(model: ENV.fetch("OLLAMA_MODEL", "gemma3"))
    @model = model
  end

  def chat(atributos:, conocimientos:, restricciones:, contexto_dinamico:, mensaje_usuario:)
    payload = {
      model: @model,
      messages: [
        { role: "system", content: prompt_base(atributos, conocimientos, restricciones) },
        { role: "system", content: contexto_dinamico.to_s },
        { role: "user", content: mensaje_usuario.to_s }
      ],
      stream: false
    }

    response = post_json(OLLAMA_URL, payload)
    content = response.dig("message", "content")
    raise "Respuesta invalida de Ollama: falta message.content" unless content.is_a?(String)

    content
  end

  private

  def prompt_base(atributos, conocimientos, restricciones)
    [
      "Atributos:",
      JSON.pretty_generate(atributos),
      "Conocimientos:",
      JSON.pretty_generate(conocimientos),
      "Restricciones:",
      JSON.pretty_generate(restricciones)
    ].join("\n")
  end

  def post_json(uri, payload)
    request = Net::HTTP::Post.new(uri)
    request["Content-Type"] = "application/json"
    request.body = JSON.generate(payload)

    http_response = Net::HTTP.start(uri.hostname, uri.port) do |http|
      http.request(request)
    end

    unless http_response.is_a?(Net::HTTPSuccess)
      raise "Error Ollama HTTP #{http_response.code}: #{http_response.body}"
    end

    JSON.parse(http_response.body)
  end
end

class OrchestratorApp < Sinatra::Base
  configure do
    set :bind, "0.0.0.0"
    set :port, 4567
    set :show_exceptions, false
    set :agent, TomNookAgent.new
    set :llm_client, LLMClient.new
  end

  before do
    content_type :json
    headers "Access-Control-Allow-Origin" => "*",
            "Access-Control-Allow-Methods" => "GET, POST, OPTIONS",
            "Access-Control-Allow-Headers" => "Content-Type, Authorization, Accept"
  end

  options "*" do
    response.headers["Allow"] = "GET, POST, OPTIONS"
    response.headers["Access-Control-Allow-Headers"] = "Content-Type, Authorization, Accept"
    response.headers["Access-Control-Allow-Origin"] = "*"
    200
  end

  post "/api/chat" do
    payload = JSON.parse(request.body.read)
    mensaje = payload["mensaje"]
    halt 400, { error: "El campo 'mensaje' es requerido y debe ser texto." }.to_json unless mensaje.is_a?(String) && !mensaje.strip.empty?

    resultado = settings.agent.evaluar_mensaje(mensaje)
    contexto = resultado[:contexto].is_a?(Array) ? resultado[:contexto].join("\n") : resultado[:contexto].to_s

    case resultado[:accion]
    when :responder_directo
      status 200
      { respuesta: contexto }.to_json
    when :llm
      respuesta_llm = settings.llm_client.chat(
        atributos: settings.agent.class::ATRIBUTOS,
        conocimientos: settings.agent.class::CONOCIMIENTOS,
        restricciones: settings.agent.class::RESTRICCIONES,
        contexto_dinamico: contexto,
        mensaje_usuario: mensaje
      )

      status 200
      { respuesta: respuesta_llm }.to_json
    else
      halt 500, { error: "Accion de agente no soportada: #{resultado[:accion].inspect}" }.to_json
    end
  rescue JSON::ParserError
    halt 400, { error: "JSON invalido." }.to_json
  end

  error StandardError do
    status 500
    { error: env["sinatra.error"].message }.to_json
  end

  run! if app_file == $PROGRAM_NAME
end
