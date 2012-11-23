// Play lowest card at every available opportunity.

#ifndef RAVAGE_H_
#define RAVAGE_H_

#include "../cc_ai_base/ai.h"
#include "../cc_ai_base/utils.h"

class RavagePlayer : public AbstractTrickTakingPlayer {
 public:
  RavagePlayer(int seat, const CardSet& hand)
      : AbstractTrickTakingPlayer(seat, hand) {}

  virtual Card LeadTrick() {
    return GetLowestCard(hand());
  }

  virtual Card FollowTrick(const vector<Card>& played_cards) {
    Card card = GetHighestLoser(hand(), played_cards);
    if (card.IsValid()) return card;
    return GetLowestWinner(hand(), played_cards);
  }

  virtual Card DiscardTrick(const vector<Card>& played_cards) {
    return GetHighestCard(hand());
  }
};

#endif  // RAVAGE_H_
