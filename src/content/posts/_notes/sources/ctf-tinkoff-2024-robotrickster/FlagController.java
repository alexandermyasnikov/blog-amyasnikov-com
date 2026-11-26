public class FlagController {
    public byte[] flag = new byte[0];
    private GameService gameService;

    public FlagController(GameService gameService2) {
        this.gameService = gameService2;
        gameService2.start(UUIDController.uuid).enqueue(new Callback<ResponseBody>() {
            public void onFailure(Call<ResponseBody> call, Throwable th) {
            }

            public void onResponse(Call<ResponseBody> call, Response<ResponseBody> response) {
                try {
                    FlagController.this.flag = response.body().bytes();
                } catch (IOException e) {
                    throw new RuntimeException(e);
                }
            }
        });
    }

    public String decrypt(byte[] bArr) {
        byte[] bArr2 = new byte[this.flag.length];
        int i = 0;
        while (true) {
            byte[] bArr3 = this.flag;
            if (i >= bArr3.length) {
                return new String(bArr2, StandardCharsets.UTF_8);
            }
            bArr2[i] = (byte) (bArr3[i] ^ bArr[i % bArr.length]);
            i++;
        }
    }
}