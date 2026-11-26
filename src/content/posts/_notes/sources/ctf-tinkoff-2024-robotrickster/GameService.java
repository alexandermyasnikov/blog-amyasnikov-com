public interface GameService {
    @GET("/accept/{uuid}/{answer}")
    Call<GameState> accept(@Path("uuid") String str, @Path("answer") Integer num);

    @GET("/next/{uuid}")
    Call<List<List<Integer>>> getNext(@Path("uuid") String str);

    @PUT("/start/{uuid}")
    Call<ResponseBody> start(@Path("uuid") String str);
}