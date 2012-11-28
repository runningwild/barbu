#include "../cc_ai_base/ai.h"
#include "lasttwo.h"
#include "ravage.h"

class MetaPlayer : public AbstractMetaPlayer {
 public:
  AbstractTrickTakingPlayer* NewRavagePlayer(int seat,
                                             const CardSet& card_set) const {
    return new RavagePlayer(seat, card_set);
  }

  AbstractTrickTakingPlayer* NewLastTwoPlayer(int seat,
                                              const CardSet& card_set) const {
    return new LastTwoPlayer(seat, card_set);
  }

  AbstractTrickTakingPlayer* NewBarbuPlayer(int seat,
                                            const CardSet& card_set) const {
    return new RavagePlayer(seat, card_set);
  }

  AbstractTrickTakingPlayer* NewKillerQueensPlayer(int seat,
                                            const CardSet& card_set) const {
    return new RavagePlayer(seat, card_set);
  }
};

int main() {
  MetaPlayer meta_player;
  meta_player.Play();
  return 0;
}
