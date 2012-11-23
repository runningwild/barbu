#ifndef AI_H_
#define AI_H_

#include <string>
#include <vector>
using namespace std;

class Card {
 public:
  Card() : suit_(0), value_(0) {}
  Card(int suit, int value) : suit_(suit), value_(value) {}

  Card(const string& s);
  static vector<Card> ListFromString(const string& s);

  int suit() const { return suit_; }
  int value() const { return value_; }

  string ToString() const;

 private:
  int suit_, value_;
};

class CardSet {
 public:
  CardSet(const vector<Card>& cards);
  CardSet(const string& s);
  static CardSet CompleteDeck();

  int GetNumCards() const { return num_cards_; }
  const vector<int>& GetValues(int suit) const;
  vector<Card> GetCards() const;
  string ToString() const;

  void RemoveCard(const Card& card);

 private:
  void Init(const vector<Card>& cards);

  int num_cards_;
  vector<vector<int> > values_by_suit_;
};

class AbstractPlayer {
 public:
  AbstractPlayer(int seat);
  virtual ~AbstractPlayer() {}

  void RunDoubling();
  virtual void PlayHand() = 0;

 protected:
  int seat() const { return seat_; }
  const vector<vector<bool> >& doubles() const { return doubles_; }

  // OVERRIDE THESE FUNCTIONS
  virtual void RecordDouble(int index, const vector<bool>& doubles);
  virtual bool ShouldDouble(int player) { return false; }

 private:
  int seat_;
  vector<vector<bool> > doubles_;
};

class AbstractTrickTakingPlayer : public AbstractPlayer {
 public:
  AbstractTrickTakingPlayer(int seat, const CardSet& hand);
  virtual ~AbstractTrickTakingPlayer() {}

  void PlayHand();

 protected:
  const CardSet& hand() const { return hand_; }
  const CardSet& remaining_cards() const { return remaining_cards_; }
  int GetWinnerIndex(const vector<Card>& played_cards) const;

  // OVERRIDE THESE FUNCTIONS
  virtual void PrepareForTrick() {}
  virtual Card LeadTrick() = 0;
  virtual Card FollowTrick(const vector<Card>& played_cards) = 0;
  virtual Card DiscardTrick(const vector<Card>& played_cards) = 0;
  virtual Card RecordTrick(const vector<Card>& played_cards,
                           int my_card_index);

 private:
  CardSet hand_, remaining_cards_;
};

class AbstractMetaPlayer {
 public:
  virtual ~AbstractMetaPlayer() {}

  void Play();

 protected:
  // OVERRIDE THESE FUNCTIONS.
  virtual AbstractTrickTakingPlayer* NewRavagePlayer(
      int seat, const CardSet& card_set) const = 0;
  virtual AbstractTrickTakingPlayer* NewLastTwoPlayer(
      int seat, const CardSet& card_set) const = 0;
  virtual AbstractTrickTakingPlayer* NewBarbuPlayer(
      int seat, const CardSet& card_set) const = 0;
};

#endif  // AI_H_
