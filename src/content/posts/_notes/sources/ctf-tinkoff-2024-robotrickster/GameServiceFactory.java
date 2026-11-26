public class GameServiceFactory {
    public GameService create() {
        return (GameService) new Retrofit.Builder()
            .baseUrl("https://t-trickster-jbi8aw9z.spbctf.ru/")
            .addConverterFactory(GsonConverterFactory.create()).build()
            .create(GameService.class);
    }
}