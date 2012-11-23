#ifndef LASTTWO_H_
#define LASTTWO_H_

#include "../cc_ai_base/ai.h"
#include "../cc_ai_base/utils.h"

class LastTwoPlayer : public AbstractTrickTakingPlayer {
 public:
  LastTwoPlayer(int seat, const CardSet& hand)
      : AbstractTrickTakingPlayer(seat, hand) {}

  virtual Card LeadTrick() {
    return GetHighestCard(hand());
  }

  virtual Card FollowTrick(const vector<Card>& played_cards) {
    return GetHighestCardInSuit(hand(), played_cards[0].suit());
  }

  virtual Card DiscardTrick(const vector<Card>& played_cards) {
    return GetHighestCard(hand());
  }
};

#endif  // LASTTWO_H_
