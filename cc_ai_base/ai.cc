#include "ai.h"

#include <algorithm>
#include <assert.h>
#include <iostream>
#include <sstream>
#include <string>
#include <vector>
using namespace std;

Card::Card(int suit, int value)
    : suit_(suit),
      value_(value) {
  assert(suit >= 0 && suit < 4);
  assert(value >= 0 && value < 13);
}

Card::Card(const string& s) {
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

vector<Card> Card::ListFromString(const string& s) {
  stringstream in(s);
  vector<Card> result;
  string temp;
  while (in >> temp) {
    result.push_back(Card(temp));
  }
  return result;
}

string Card::ToString() const {
  string result = "??";
  if (!IsValid()) return result;
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

CardSet::CardSet(const vector<Card>& cards) {
  Init(cards);
}

CardSet::CardSet(const string& s) {
  Init(Card::ListFromString(s));
}

CardSet CardSet::CompleteDeck() {
  vector<Card> cards;
  for (int suit = 0; suit < 4; ++suit)
    for (int value = 0; value < 13; ++value)
      cards.push_back(Card(suit, value));
  return CardSet(cards);
}

const vector<int>& CardSet::GetValues(int suit) const {
  assert(suit >= 0 && suit < 4);
  return values_by_suit_[suit];
}

vector<Card> CardSet::GetCards() const {
  vector<Card> result;
  for (int suit = 0; suit < 4; ++suit) {
    for (int i = 0; i < values_by_suit_[suit].size(); ++i) {
      result.push_back(Card(suit, values_by_suit_[suit][i]));
    }
  }
  return result;
}

string CardSet::ToString() const {
  ostringstream out;
  vector<Card> cards = GetCards();
  for (int i = 0; i < cards.size(); ++i) {
    if (i != 0) out << " ";
    out << cards[i].ToString();
  }
  return out.str();
}

void CardSet::RemoveCard(const Card& card) {
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

void CardSet::Init(const vector<Card>& cards) {
  num_cards_ = cards.size();
  values_by_suit_.resize(4);
  for (int i = 0; i < cards.size(); ++i) {
    values_by_suit_[cards[i].suit()].push_back(cards[i].value());
  }
  for (int i = 0; i < 4; ++i) {
    sort(values_by_suit_[i].begin(), values_by_suit_[i].end());
  }
}

AbstractPlayer::AbstractPlayer(int seat)
  : seat_(seat),
    doubles_(4) {}

void AbstractPlayer::RunDoubling() {
  string doubling_header;
  getline(cin, doubling_header);
  assert(doubling_header == "DOUBLING");

  for (int i = 0; i < 4; i++) {
    int index = (i == 3 ? 0 : i + 1);
    string doubling, part;
    getline(cin, doubling);
    stringstream in(doubling);
    vector<string> parts;
    while (in >> part) parts.push_back(part);

    // Our chance to double.
    if (parts.size() == 1) {
      ostringstream out;
      assert(parts[0] == "DOUBLE");
      vector<int> to_double;
      for (int j = 0; j < 4; ++j) {
        if (j == index) continue;
        if (i == 3 && !doubles()[j][index]) continue;
        if (ShouldDouble(j))
          out << " " << j;
      }
      string s = out.str();
      if (s.empty())
        cout << "" << endl;
      else
        cout << s.substr(1) << endl;
    }

    // Someone else is doubling.
    else {
      vector<bool> to_double(4, false);
      assert(parts[1] == "DOUBLE");
      for (int j = 2; j < parts.size(); ++j)
        to_double[atoi(parts[j].c_str())] = true;
      RecordDouble(index, to_double);
    }
  }  
}

void AbstractPlayer::RecordDouble(
    int index, const vector<bool>& doubles) {
  doubles_[index] = doubles;
}

AbstractTrickTakingPlayer::AbstractTrickTakingPlayer(
    int seat,
    const CardSet& hand)
    : AbstractPlayer(seat),
      hand_(hand),
      remaining_cards_(CardSet::CompleteDeck()) {
  vector<Card> cards = hand.GetCards();
  for (int i = 0; i < cards.size(); ++i)
    remaining_cards_.RemoveCard(cards[i]);
}

void AbstractTrickTakingPlayer::PlayHand() {
  string line;
  while (true) {
    getline(cin, line);
    if (line != "TRICK") {
      break;
    }

    vector<Card> full_trick;
    int position;
    for (int i = 0; i < 4; i++) {
      getline(cin, line);
      if (line == "PLAY") {
	Card play;
	position = i;
	PrepareForTrick();
	if (full_trick.empty()) {
	  play = LeadTrick();
	} else if (hand_.GetValues(full_trick[0].suit()).size()) {
	  play = FollowTrick(full_trick);
	} else {
	  play = DiscardTrick(full_trick);
	}
	cout << play.ToString() << endl;
	full_trick.push_back(play);
      } else {
	int seat;
	string card;
	stringstream ss(line);
	ss >> seat >> card;
	full_trick.push_back(card);
      }
    }
    RecordTrick(full_trick, position);
  }
  assert (line.substr(0, 3) == "END");
}

Card AbstractTrickTakingPlayer::RecordTrick(const vector<Card>& played_cards,
					    int my_card_index) {
  for (int i = 0; i < played_cards.size(); ++i) {
    if (i == my_card_index) {
      hand_.RemoveCard(played_cards[i]);
    } else {
      remaining_cards_.RemoveCard(played_cards[i]);
    }
  }
}

void AbstractMetaPlayer::Play() {
  while (true) {
    string player_line, temp;
    int seat;
    if (!getline(cin, player_line)) return;
    stringstream in(player_line);
    in >> temp >> seat;

    string hand;
    getline(cin, hand);
    CardSet card_set = CardSet(hand);

    string game;
    getline(cin, game);

    AbstractTrickTakingPlayer* player = NULL;
    if (game == "RAVAGE") {
      player = NewRavagePlayer(seat, card_set);
    } else if (game == "LASTTWO") {
      player = NewLastTwoPlayer(seat, card_set);
    } else if (game == "BARBU") {
      player = NewRavagePlayer(seat, card_set);
    } else if (game == "KILLERQUEENS") {
      player = NewKillerQueensPlayer(seat, card_set);
    } else {
      assert(false);
    }

    player->RunDoubling();
    player->PlayHand();
    delete player;
  }
}
