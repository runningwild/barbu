#! /usr/bin/python

import sys


class Card(object):
  _value_map = {'2': 0, '3': 1, '4': 2, '5': 3, '6': 4, '7': 5, '8': 6,
                '9': 7, 't': 8, 'j': 9, 'q': 10, 'k': 11, 'a': 12}
  _suit_map= {'c': 0, 'd': 1, 'h': 2, 's': 3}
  _values = ['2', '3', '4', '5', '6', '7', '8', '9', 't', 'j', 'q', 'k', 'a']
  _suits = ['c', 'd', 'h', 's']

  @staticmethod
  def FromString(s):
    return Card(Card._suit_map[s[1]], Card._value_map[s[0]])

  @staticmethod
  def ListFromString(s):
    return [Card.FromString(word) for word in s.split()]

  def __init__(self, suit, value):
    self.value = value
    self.suit = suit

  def __str__(self):
    return Card._values[self.value] + Card._suits[self.suit]

  def __repr__(self):
    return self.__str__()


class CardSet(object):
  @staticmethod
  def CompleteDeck():
    cards = []
    for suit in xrange(4):
      for value in xrange(13):
        cards.append(Card(suit, value))
    return CardSet(cards)

  def __init__(self, cards):
    self._values_by_suit = [[], [], [], []]
    self._num_cards = 0
    for card in cards:
      self._values_by_suit[card.suit].append(card.value)
      self._num_cards += 1
    for values in self._values_by_suit:
      values.sort()
    
  def RemoveCard(self, c):
    values = self._values_by_suit[c.suit]
    index = None
    for i, value in enumerate(values):
      if value == c.value:
        index = i
        break
    assert index != None
    del values[index]
    self._num_cards -= 1

  def GetNumCards(self):
    return self._num_cards

  def GetValues(self, suit):
    return self._values_by_suit[suit]

  def __str__(self):
    result = []
    for suit, values in enumerate(self._values_by_suit):
      for value in values:
        result.append(str(Card(suit, value)))
    return ' '.join(result)

# Functions that must be implemented:
#  - GetHandResultDescription(self):
#      Returns a description of how the hand went, after it is over.
#  - _ChooseLead(self, suit):
#      Returns a pair of (card, goodness) to lead in the given suit. It is
#      guaranteed that at least one valid play exists. The returned goodness
#      must be positive.
#  - _ChooseFollow(self, suit, played_values, discards):
#      Returns the value of a card to play in the given suit. played_values is
#      the list of values already played in the suit (in order). discards is the
#      list of cards played from other suits (in order).
#  - _ChooseDiscard(self, discard_suit, suit, played_values, discards):
#      Returns a pair of (card, goodness) to discard in the given discard suit.
#      The rest of the parameters are as in _ChooseFollow.
class AbstractTrickTakingPlayer(object):
  def __init__(self, hand_cards):
    self.hand = CardSet(hand_cards)
    self.remaining_cards = CardSet.CompleteDeck()
    for card in hand_cards:
      self.remaining_cards.RemoveCard(card)

  def IsDone(self):
    return self.hand.GetNumCards() == 0

  def LeadTrick(self):
    best_card, best_goodness = None, 0
    for suit in xrange(4):
      if self.hand.GetValues(suit):
        value, goodness = self._ChooseLead(suit)
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
    if self.hand.GetValues(suit):
      return Card(suit, self._ChooseFollow(suit, played_values, discards))
    else:
      best_discard, best_goodness = None, 0
      for discard_suit in xrange(4):
        if self.hand.GetValues(discard_suit):
          value, goodness = self._ChooseDiscard(discard_suit, suit,
                                                played_values, discards)
          if goodness > best_goodness:
            best_discard = Card(discard_suit, value)
            best_goodness = goodness
      assert best_discard is not None
      return best_discard

  def RecordTrick(self, cards, my_card_index):
    winner = cards[0]
    for i, card in enumerate(cards):
      if i == my_card_index:
        self.hand.RemoveCard(card)
      else:
        self.remaining_cards.RemoveCard(card)
      if card.suit == winner.suit and card.value > winner.value:
        winner = card
    return winner


def RunTrickTakingPlayer(player, print_history=False):
  history = []
  history.append('HAND=[%s]' % player.hand)
  trick_number = 1
  while True:
    if player.IsDone(): break
    line = sys.stdin.readline()
    if line == 'quit': break
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
    history.append('#%d: TRICK=%s PLAYED=%s'
                   % (trick_number, whole_trick, play))
    trick_number += 1
  if print_history:
    history.append('HAND RESULT: %s' % player.GetHandResultDescription())
    print >> sys.stderr, '\n'.join(history)
