/*Data structure containing LED information*/
typedef struct {
  int addr;
  bool updir;
  int prior;
  int pace;
  int eff;
  int lumA;
  int colorA;
  int lumB;
  int colorB;
} DataStream;
