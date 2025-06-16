# to create a vocabulary of keywords

class Vocabulary:
    def __init__(self):
        self.vocab = []
        self.seen = set()

    def createVocab(self, keywords):
        self.seen = set()
        self.vocab = []
        for filename, keywords in keywords.items():
            # print(f"\n== FILE: {filename} ==")
            for kw, score in keywords:
                if kw not in self.seen:
                    self.seen.add(kw)
                    self.vocab.append(kw)
        return self.vocab
    
    def addToVocab(self, keywords):
        vocab = []
        for filename, kw_list in keywords.items():
            for kw, score in kw_list:
                if kw not in self.seen:
                    self.seen.add(kw)
                    vocab.append(kw)
        self.vocab.append(vocab)
        return self.vocab