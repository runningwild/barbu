#ifndef RAVAGE_H_
#define RAVAGE_H_

#include "../cc_ai_base/ai.h"
#include "../cc_ai_base/utils.h"
#include <assert.h>

class RavagePlayer : public AbstractTrickTakingPlayer {
public:
  RavagePlayer(int seat, const CardSet& hand)
    : AbstractTrickTakingPlayer(seat, hand),
      taken_by_suit_(4, 0) {}

  virtual bool ShouldDouble(int player) { return true; }

  virtual void PrepareForTrick() {
    RecalculateBadness(&card_badness_, &suit_badness_);
  }

  virtual Card LeadTrick() {
    Card best_card;
    double best_goodness = -9999;
    for (int suit = 0; suit < 4; ++suit) {
      const vector<int>& my_values = hand().GetValues(suit);
      const vector<int>& remaining_values = remaining_cards().GetValues(suit);
      if (my_values.empty()) continue;

      double goodness = 1.0 * remaining_values.size() / my_values.size();
      if (goodness > best_goodness) {
	best_card = Card(suit, my_values[0]);
	if (my_values.size() > 1 &&
	    taken_by_suit_[suit] < 4 &&
	    remaining_values.size() >= 4) {
	  best_card = Card(suit, my_values[1]);
	}
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
    const vector<int>& other_values = remaining_cards().GetValues(suit);

    int highest_duck_index = -1, lowest_win_index = -1;
    for (int i = 0; i < my_values.size(); ++i) {
      if (my_values[i] < winner_value)
	highest_duck_index = i;
      else if (lowest_win_index == -1)
	lowest_win_index = i;
    }

    if (highest_duck_index != -1) {
      return Card(suit, my_values[highest_duck_index]);
    }

    assert(lowest_win_index != -1);
    if (taken_by_suit_[suit] < 4 || played_cards.size() >= 2) {
      return Card(suit, my_values.back());
    }
    return Card(suit, my_values[lowest_win_index]);
  }

  virtual Card DiscardTrick(const vector<Card>& played_cards) {
    Card best_card;
    double best_goodness = -9999;
    for (int suit = 0; suit < 4; ++suit) {
      const vector<int>& my_values = hand().GetValues(suit);
      const vector<int>& other_values = remaining_cards().GetValues(suit);
      if (my_values.empty()) continue;
      double goodness = 1000 - my_values.size();
      if (other_values.empty()) goodness = 1;
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
  void RecalculateBadness(vector<vector<double> >* card_badness,
			  vector<double>* suit_badness) const {
    card_badness->resize(4);
    suit_badness->resize(4, 0);
    for (int suit = 0; suit < 4; ++suit) {
      const vector<int>& my_values = hand().GetValues(suit);
      const vector<int>& other_values = remaining_cards().GetValues(suit);
      int j = 0;
      for (int i = 0; i < my_values.size(); ++i) {
	while (j < other_values.size() && other_values[j] < my_values[i])
	  ++j;
	int rank = i + j;
	int base_rank = 3 + 4 * i;
	double card_cost = 0.5 + (rank - base_rank) * 0.25;
	card_cost = min(card_cost, 1.0);
	card_cost = max(card_cost, 0.0);
	(*card_badness)[suit].push_back(card_cost);
	(*suit_badness)[suit] += card_cost;
      }
    }
  }

  vector<vector<double> > card_badness_;
  vector<double> suit_badness_;
  vector<int> taken_by_suit_;
};

#endif  // RAVAGE_H_
