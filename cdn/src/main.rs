use actix_web::{get, post, delete, web::{self, PayloadConfig}, App, HttpResponse, HttpServer. Responder, Error, HttpRequest};

#[get("/")]
async fn greetings() -> impl Responder {
    HttpResponse::Ok().body("greetings user of the internet...")
}

#[actix_web::main] 
async fn main() -> std::io::Result<()> {
  HttpServer::new(|| {
    App::new()
      .service(greetings)
      .app_data(PayloadConfig::new(20_000_000_000))
  })
  .bind("127.0.0.1:8080")?
  .run()
  .await
}