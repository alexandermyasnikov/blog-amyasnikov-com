public class GameViewModel extends AndroidViewModel {
    private MutableLiveData<String> flag = new MutableLiveData<>("");
    private FlagController flagController = new FlagController(this.gameService);

    // ...

    public void startGame() {
        this.blockPlayButton.setValue(false);
        this.score.setValue(0);
    }

    public List<Integer> getCurrentPosition() {
        return this.currentPosition.getValue();
    }

    public void getNextTransposition(Callback<List<List<Integer>>> callback) {
        this.blockPlayButton.setValue(true);
        this.gameService.getNext(UUIDController.uuid).enqueue(callback);
    }

    public void clickByCup(int i) {
        if (this.canTryAnswer) {
            stopChooseAnswer();
            this.gameService.accept(UUIDController.uuid, Integer.valueOf(this.currentPosition.getValue().lastIndexOf(Integer.valueOf(i)))).enqueue(new Callback<GameState>() {
                public void onResponse(Call<GameState> call, Response<GameState> response) {
                    GameViewModel.this.blockPlayButton.setValue(false);
                    GameViewModel.this.score.setValue(Integer.valueOf(response.body().getScore()));
                    GameViewModel.this.update(response.body().getKey());
                }

                public void onFailure(Call<GameState> call, Throwable th) {
                    GameViewModel.this.blockPlayButton.setValue(false);
                    GameViewModel.this.update(new byte[0]);
                }
            });
        }
    }

    public void update(byte[] bArr) {
        if (bArr != null && bArr.length != 0) {
            this.blockPlayButton.setValue(true);
            this.flag.setValue("The android swindler is defeated!\nHis last words:\n" + this.flagController.decrypt(bArr));
        }
    }

    // ...
}