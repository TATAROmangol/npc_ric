"""add institution_id to Template

Revision ID: 0d4f9158728e
Revises: 
Create Date: 2025-05-23 12:00:55.325974
"""

from typing import Sequence, Union
from alembic import op
import sqlalchemy as sa

# revision identifiers, used by Alembic.
revision: str = '0d4f9158728e'
down_revision: Union[str, None] = None
branch_labels: Union[str, Sequence[str], None] = None
depends_on: Union[str, Sequence[str], None] = None


def upgrade() -> None:
    op.add_column('templates', sa.Column('institution_id', sa.Integer(), nullable=True))
    op.drop_constraint('templates_name_key', 'templates', type_='unique')
    op.create_unique_constraint(None, 'templates', ['institution_id'])


def downgrade() -> None:
    op.drop_constraint(None, 'templates', type_='unique')
    op.create_unique_constraint('templates_name_key', 'templates', ['name'])
    op.drop_column('templates', 'institution_id')
