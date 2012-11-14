#ifndef RAVAGE_H_
#define RAVAGE_H_

#include "base.h"

class RavagePlayer : public AbstractTrickTakingPlayer {
public:
  RavagePlayer(const CardSet& hand)
    : AbstractTrickTakingPlayer(hand),
      taken_by_suit_(4, 0) {}

  virtual Card LeadTrick() const {
    Card best_card;
    double best_goodness = -9999;
    for (int suit = 0; suit < 4; ++suit) {
      const vector<int>& my_values = hand().GetValues(suit);
      const vector<int>& remaining_values = remaining_cards().GetValues(suit);
      if (my_values.empty()) continue;
      double goodness = 1.0 * remaining_values.size() / my_values.size();
      //double goodness = 10000 - GetSuitBadness(suit, false);
      if (goodness > best_goodness) {
	best_card = Card(suit, my_values[0]);
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
    const vector<int>& other_values = remaining_cards().GetValues(suit);

    /*int adjusted_winner_index = min(3, int(other_values.size()) - 1);
    int adjusted_winner_value = winner_value;
    if (adjusted_winner_index >= 0) {
      adjusted_winner_value = max(adjusted_winner_value,
				  other_values[adjusted_winner_index]);
				  }*/

    int highest_duck_index = -1, lowest_win_index = -1;
    //int adjusted_highest_duck_index = -1, adjusted_lowest_win_index = -1;
    for (int i = 0; i < my_values.size(); ++i) {
      if (my_values[i] < winner_value)
	highest_duck_index = i;
      else if (lowest_win_index == -1)
	lowest_win_index = i;
      /*if (my_values[i] < adjusted_winner_value)
	adjusted_highest_duck_index = i;
      else if (adjusted_lowest_win_index == -1)
      adjusted_lowest_win_index = i;*/
    }

    RecalculateBadness();
    if (highest_duck_index != -1 &&
	card_badness_[suit][highest_duck_index] > 0) {
      return Card(suit, my_values[highest_duck_index]);
    }
    if (lowest_win_index != -1) {
      for (int i = lowest_win_index; i < card_badness_[suit].size(); ++i)
	if (card_badness_[suit][i] >= 0.6)
	  return Card(suit, my_values.back());
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

  virtual Card DiscardTrick(const vector<Card>& played_cards) const {
    Card best_card;
    double best_goodness = -9999;
    for (int suit = 0; suit < 4; ++suit) {
      const vector<int>& my_values = hand().GetValues(suit);
      const vector<int>& other_values = remaining_cards().GetValues(suit);
      if (my_values.empty()) continue;
      //double goodness = 1.0 * remaining_values.size() / my_values.size();
      double goodness = 1000 - my_values.size();
      if (other_values.empty()) goodness = 1;
      //double delta = GetSuitBadness(suit, false) - GetSuitBadness(suit, true);
      //double goodness = 1 + delta;
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
  void RecalculateBadness() const {
    RecalculateBadness(&card_badness_, &suit_badness_);
  }
  mutable vector<vector<double> > card_badness_;
  mutable vector<double> suit_badness_;

  void RecalculateBadness(vector<vector<double> >* card_badness,
			  vector<double>* suit_badness) const {
    card_badness->resize(4);
    suit_badness->resize(4, 0);
    for (int suit = 0; suit < 4; ++suit) {
      const vector<int>& my_values = hand().GetValues(suit);
      const vector<int>& other_values = remaining_cards().GetValues(suit);

      for (int i = 0; i < my_values.size(); ++i) {
	int rank = i;
	for (int j = 0; j < other_values.size()
	       && other_values[j] < my_values[i]; ++j) {
	  ++rank;
	}
	int base_rank = 3 + 4 * i;
	double card_cost = 0.5 + (rank - base_rank) * 0.25;
	card_cost = min(card_cost, 1.0);
	card_cost = max(card_cost, 0.0);
	(*card_badness)[suit].push_back(card_cost);
	(*suit_badness)[suit] += card_cost;
      }
    }
  }

  double GetSuitBadness(int suit, bool drop_highest) const {
    vector<int> my_values = hand().GetValues(suit);
    if (drop_highest) my_values.pop_back();
    const vector<int>& other_values = remaining_cards().GetValues(suit);
    int my_pos = 0, other_pos = 0;

    double result = 0;
    while (my_pos < my_values.size() &&
	   other_pos < other_values.size()) {
      int beaters = 0, duckers = 0;
      for (int i = 0; i < 3; ++i) {
	if (other_values[other_pos] < my_values[my_pos])
	  ++duckers;
	else
	  ++beaters;
	++other_pos;
	if (other_pos >= other_values.size()) break;
      }
      result += duckers * 1.0 / (duckers + beaters);
      my_pos += 1;
    }
    result -= (other_values.size() - other_pos) * 0.25;
    result += taken_by_suit_[suit] * 0.25;
    result += 5;
    return result * result;
  }

  vector<int> taken_by_suit_;
};

#endif  // RAVAGE_H_
