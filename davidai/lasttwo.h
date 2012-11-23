#ifndef LASTTWO_H_
#define LASTTWO_H_

#include "../cc_ai_base/ai.h"
#include <assert.h>

class LastTwoPlayer : public AbstractTrickTakingPlayer {
public:
  LastTwoPlayer(int seat, const CardSet& hand)
    : AbstractTrickTakingPlayer(seat, hand) {}

  virtual Card LeadTrick() {
    Card best_card;
    double best_goodness = -9999;
    for (int suit = 0; suit < 4; ++suit) {
      const vector<int>& my_values = hand().GetValues(suit);
      const vector<int>& remaining_values = remaining_cards().GetValues(suit);
      if (my_values.empty()) continue;
      double goodness = 0;
      goodness = my_values.back();
      if (goodness > best_goodness) {
	best_card = Card(suit, my_values.back());
	best_goodness = goodness;
      }
    }
    assert(best_goodness >= 0);
    return best_card;
  }

  virtual Card FollowTrick(const vector<Card>& played_cards) {
    int suit = played_cards[0].suit();
    int winner_value = played_cards[GetWinnerIndex(played_cards)].value();
    const vector<int>& my_values = hand().GetValues(suit);
    return Card(suit, my_values.back());
  }

  virtual Card DiscardTrick(const vector<Card>& played_cards) {
    Card best_card;
    double best_goodness = -9999;
    for (int suit = 0; suit < 4; ++suit) {
      const vector<int>& my_values = hand().GetValues(suit);
      if (my_values.empty()) continue;
      double goodness = my_values.back();
      if (goodness >= best_goodness) {
	best_card = Card(suit, my_values.back());
	best_goodness = goodness;
      }
    }
    assert(best_goodness >= 0);
    return best_card;
  }
};

#endif  // LASTTWO_H_
