#! /usr/bin/python

import sys


class Card:
  value_map_ = {'2': 0, '3': 1, '4': 2, '5': 3, '6': 4, '7': 5, '8': 6,
                '9': 7, 't': 8, 'j': 9, 'q': 10, 'k': 11, 'a': 12}
  suit_map_ = {'c': 0, 'd': 1, 'h': 2, 's': 3}
  values_ = ['2', '3', '4', '5', '6', '7', '8', '9', 't', 'j', 'q', 'k', 'a']
  suits_ = ['c', 'd', 'h', 's']

  @staticmethod
  def FromString(s):
    return Card(Card.suit_map_[s[1]], Card.value_map_[s[0]])

  @staticmethod
  def ListFromString(s):
    return [Card.FromString(word) for word in s.split()]

  def __init__(self, suit, value):
    self.value = value
    self.suit = suit

  def __str__(self):
    return Card.values_[self.value] + Card.suits_[self.suit]

  def __repr__(self):
    return self.__str__()


class CardSet:
  @staticmethod
  def CompleteDeck():
    cards = []
    for suit in xrange(4):
      for value in xrange(13):
        cards.append(Card(suit, value))
    return CardSet(cards)

  def __init__(self, cards):
    self.values_by_suit_ = [[], [], [], []]
    self.num_cards_ = 0
    for card in cards:
      self.values_by_suit_[card.suit].append(card.value)
      self.num_cards_ += 1
    for values in self.values_by_suit_:
      values.sort()
    
  def RemoveCard(self, c):
    values = self.values_by_suit_[c.suit]
    index = None
    for i, value in enumerate(values):
      if value == c.value:
        index = i
        break
    assert index != None
    del values[index]
    self.num_cards_ -= 1

  def GetNumCards(self):
    return self.num_cards_

  def GetValues(self, suit):
    return self.values_by_suit_[suit]

  def __str__(self):
    result = []
    for suit, values in enumerate(self.values_by_suit_):
      for value in values:
        result.append(str(Card(suit, value)))
    return ' '.join(result)


class RavagePlayer:
  def __init__(self, hand_cards):
    self.remaining_cards_ = CardSet.CompleteDeck()
    self.hand_ = CardSet(hand_cards)
    self.taken_by_suit_ = [0, 0, 0, 0]
    for card in hand_cards:
      self.remaining_cards_.RemoveCard(card)

  def LeadTrick(self):
    best_card, best_goodness = None, 0
    for suit in xrange(4):
      if self.hand_.GetValues(suit):
        value, goodness = self._EvaluateLead(suit)
        if goodness > best_goodness:
          best_card = Card(suit, value)
          best_goodness = goodness
    assert best_card is not None
    return best_card

  def FollowTrick(self, played_cards):
    suit = played_cards[0].suit
    played_values = []
    discards = []
    for card in played_cards:
      if card.suit == suit:
        played_values.append(card.value)
      else:
        discards.append(card)
    if self.hand_.GetValues(suit):
      return Card(suit, self._ChooseFollow(suit, played_values, discards))
    else:
      best_discard, best_goodness = None, 0
      for discard_suit in xrange(4):
        if self.hand_.GetValues(discard_suit):
          value, goodness = self._ChooseDiscard(discard_suit, suit,
                                                played_values, discards)
          if goodness > best_goodness:
            best_discard = Card(discard_suit, value)
            best_goodness = goodness
      assert best_discard is not None
      return best_discard

  def RecordTrick(self, cards, my_card_index):
    cards_by_suit = [0, 0, 0, 0]
    winner = cards[0]
    for i, card in enumerate(cards):
      cards_by_suit[card.suit] += 1
      if i == my_card_index:
        self.hand_.RemoveCard(card)
      else:
        self.remaining_cards_.RemoveCard(card)
      if card.suit == winner.suit and card.value > winner.value:
        winner = card
    if winner == cards[my_card_index]:
      for suit, number in enumerate(cards_by_suit):
        self.taken_by_suit_[suit] += number

  def _EvaluateLead(self, suit):
    my_values = self.hand_.GetValues(suit)
    return (my_values[0], 20 - my_values[0])

  def _ChooseFollow(self, suit, played_values, discards):
    num_played = len(played_values) + len(discards)
    highest_played = max(played_values)
    my_values = self.hand_.GetValues(suit)
    for (i, value) in enumerate(my_values):
      if i + 1 == len(my_values) or my_values[i + 1] > highest_played:
        return my_values[i]

  def _ChooseDiscard(self, discard_suit, suit, played_values, discards):
    my_values = self.hand_.GetValues(discard_suit)
    return (my_values[-1], my_values[-1] + 1)


hand_cards = Card.ListFromString(sys.stdin.readline())
player = RavagePlayer(hand_cards)
for _ in range(13):
while True:
  line = sys.stdin.readline()
  if line == 'quit' or not line: break
  new_trick = Card.ListFromString(line)
  if new_trick:
    play = player.FollowTrick(new_trick)
  else:
    play = player.LeadTrick()
  print str(play)

  sys.stdout.flush()
  rest_of_trick = Card.ListFromString(sys.stdin.readline())
  whole_trick = new_trick + [play] + rest_of_trick
  player.RecordTrick(whole_trick, len(new_trick))
