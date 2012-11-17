#ifndef LASTTWO_H_
#define LASTTWO_H_

#include "base.h"

class LastTwoPlayer : public AbstractTrickTakingPlayer {
public:
  LastTwoPlayer(const CardSet& hand)
    : AbstractTrickTakingPlayer(hand),
      taken_by_suit_(4, 0) {}

  virtual Card LeadTrick() const {
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

  virtual Card FollowTrick(const vector<Card>& played_cards) const {
    int suit = played_cards[0].suit();
    int winner_value = played_cards[GetWinnerIndex(played_cards)].value();
    const vector<int>& my_values = hand().GetValues(suit);
    return Card(suit, my_values.back());
  }

  virtual Card DiscardTrick(const vector<Card>& played_cards) const {
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

  virtual Card RecordTrick(const vector<Card>& played_cards,
			   int my_card_index) {
    if (my_card_index == GetWinnerIndex(played_cards)) {
      for (int i = 0; i < played_cards.size(); ++i) {
	taken_by_suit_[played_cards[i].suit()] += 1;
      }
    }
    AbstractTrickTakingPlayer::RecordTrick(played_cards, my_card_index);
  }

private:
  vector<int> taken_by_suit_;
};

#endif  // LASTTWO_H_
