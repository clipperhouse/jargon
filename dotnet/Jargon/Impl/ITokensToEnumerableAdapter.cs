using System;
using System.Collections;
using System.Collections.Generic;

namespace Jargon.Impl
{
    internal sealed class ITokensToEnumerableAdapter: IEnumerable<Token>
    {
        private sealed class Enumerator: IEnumerator<Token>
        {
            public Token Current { get; private set; }

            object IEnumerator.Current => Current;

            private ITokens Inner;
            public Enumerator(ITokens inner)
            {
                Inner = inner;
            }

            public bool MoveNext()
            {
                var tok = Inner.Next();
                if (tok == null) return false;

                Current = tok.Value;
                return true;
            }

            public void Dispose()
            {
                Inner?.Dispose();
                Inner = null;

            }

            public void Reset()
            {
                throw new NotImplementedException(nameof(Reset));
            }
            
        }

        private ITokens Inner;
        public ITokensToEnumerableAdapter(ITokens inner)
        {
            Inner = inner;
        }
        
        public IEnumerator<Token> GetEnumerator()
        {
            var ret = Inner;
            Inner = null;
            if (ret == null) throw new InvalidOperationException("Can only be enumerated once");

            return new Enumerator(ret);
        }

        IEnumerator IEnumerable.GetEnumerator() => GetEnumerator();
    }
}
