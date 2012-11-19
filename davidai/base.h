#ifndef BASE_H_
#define BASE_H_

#include <algorithm>
#include <assert.h>
#include <iostream>
#include <sstream>
#include <string>
#include <vector>
using namespace std;

class Card {
 public:
  Card() : suit_(0), value_(0) {}
  Card(int suit, int value) : suit_(suit), value_(value) {}

  Card(const string& s) {
    if (s[0] >= '2' && s[0] <= '9') value_ = s[0] - '2';
    else if (s[0] == 't') value_ = 8;
    else if (s[0] == 'j') value_ = 9;
    else if (s[0] == 'q') value_ = 10;
    else if (s[0] == 'k') value_ = 11;
    else if (s[0] == 'a') value_ = 12;
    else assert(false);
    if (s[1] == 'c') suit_ = 0;
    else if (s[1] == 'd') suit_ = 1;
    else if (s[1] == 'h') suit_ = 2;
    else if (s[1] == 's') suit_ = 3;
    else assert(false);
  }

  static vector<Card> ListFromString(const string& s) {
    stringstream in(s);
    vector<Card> result;
    string temp;
    while (in >> temp) {
      result.push_back(Card(temp));
    }
    return result;
  }

  int suit() const { return suit_; }
  int value() const { return value_; }

  string ToString() const {
    string result = "??";
    if (value_ >= 0 && value_ <= 7) result[0] = '2' + value_;
    else if (value_ == 8) result[0] = 't';
    else if (value_ == 9) result[0] = 'j';
    else if (value_ == 10) result[0] = 'q';
    else if (value_ == 11) result[0] = 'k';
    else if (value_ == 12) result[0] = 'a';
    else assert(false);
    if (suit_ == 0) result[1] = 'c';
    else if (suit_ == 1) result[1] = 'd';
    else if (suit_ == 2) result[1] = 'h';
    else if (suit_ == 3) result[1] = 's';
    else assert(false);
    return result;
  }

 private:
  int suit_, value_;
};

class CardSet {
 public:
  CardSet(const vector<Card>& cards) {
    Init(cards);
  }

  CardSet(const string& s) {
    Init(Card::ListFromString(s));
  }

  static CardSet CompleteDeck() {
    vector<Card> cards;
    for (int suit = 0; suit < 4; ++suit)
      for (int value = 0; value < 13; ++value)
	cards.push_back(Card(suit, value));
    return CardSet(cards);
  }

  int GetNumCards() const { return num_cards_; }
  const vector<int>& GetValues(int suit) const {
    assert(suit >= 0 && suit < 4);
    return values_by_suit_[suit];
  }

  vector<Card> GetCards() const {
    vector<Card> result;
    for (int suit = 0; suit < 4; ++suit) {
      for (int i = 0; i < values_by_suit_[suit].size(); ++i) {
	result.push_back(Card(suit, values_by_suit_[suit][i]));
      }
    }
    return result;
  }

  void RemoveCard(const Card& card) {
    vector<int>& values = values_by_suit_[card.suit()];
    int pos = 0;
    while (pos < values.size() && values[pos] != card.value()) {
      ++pos;
    }
    assert(pos < values.size());
    while (pos + 1 < values.size()) {
      values[pos] = values[pos + 1];
      ++pos;
    }
    values.pop_back();
    --num_cards_;
  }

  string ToString() const {
    ostringstream out;
    vector<Card> cards = GetCards();
    for (int i = 0; i < cards.size(); ++i) {
      if (i != 0) out << " ";
      out << cards[i].ToString();
    }
    return out.str();
  }

 private:
  void Init(const vector<Card>& cards) {
    num_cards_ = cards.size();
    values_by_suit_.resize(4);
    for (int i = 0; i < cards.size(); ++i) {
      values_by_suit_[cards[i].suit()].push_back(cards[i].value());
    }
    for (int i = 0; i < 4; ++i) {
      sort(values_by_suit_[i].begin(), values_by_suit_[i].end());
    }
  }

  int num_cards_;
  vector<vector<int> > values_by_suit_;
};

class AbstractTrickTakingPlayer {
 public:
  AbstractTrickTakingPlayer(int seat, const CardSet& hand)
    : seat_(seat), hand_(hand),
    remaining_cards_(CardSet::CompleteDeck()) {
    vector<Card> cards = hand.GetCards();
    for (int i = 0; i < cards.size(); ++i) {
      remaining_cards_.RemoveCard(cards[i]);
    }
  }

  int seat() const { return seat_; }
  const CardSet& hand() const { return hand_; }
  const CardSet& remaining_cards() const { return remaining_cards_; }
  bool IsDone() const { return hand_.GetNumCards() == 0; }

  // OVERRIDE THESE FUNCTIONS
  virtual void PrepareForTrick() {}
  virtual Card LeadTrick() const = 0;
  virtual Card FollowTrick(const vector<Card>& played_cards) const = 0;
  virtual Card DiscardTrick(const vector<Card>& played_cards) const = 0;
  virtual Card RecordTrick(const vector<Card>& played_cards,
			   int my_card_index) {
    for (int i = 0; i < played_cards.size(); ++i) {
      if (i == my_card_index) {
	hand_.RemoveCard(played_cards[i]);
      } else {
	remaining_cards_.RemoveCard(played_cards[i]);
      }
    }
  }

 protected:
  int GetWinnerIndex(const vector<Card>& played_cards) const {
    assert(!played_cards.empty());
    int index = 0;
    for (int i = 1; i < played_cards.size(); ++i)
      if (played_cards[i].suit() == played_cards[index].suit() &&
	  played_cards[i].value() > played_cards[index].value())
	index = i;
    return index;
  }

 private:
  int seat_;
  CardSet hand_, remaining_cards_;
};

#endif  // BASE_H_
